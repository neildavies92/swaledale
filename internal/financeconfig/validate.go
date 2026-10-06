package financeconfig

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/neildavies/swaledale/internal/domain"
)

var keyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

func addKey(seen map[Key]bool, key Key, path string) error {
	if !keyPattern.MatchString(string(key)) {
		return fmt.Errorf("%s: key must match [a-z][a-z0-9_]{0,63}", path)
	}
	if seen[key] {
		return fmt.Errorf("%s: duplicate key %q", path, key)
	}
	seen[key] = true
	return nil
}
func reference(seen map[Key]bool, key Key, path string) error {
	if !seen[key] {
		return fmt.Errorf("%s: unknown reference %q", path, key)
	}
	return nil
}

// Validate checks configuration references and canonical financial invariants.
// Validation is deterministic in file order and does not mutate the definitions.
func (c Config) Validate() error {
	if c.Version != Version {
		return fmt.Errorf("version: expected %d", Version)
	}
	if len(c.Members) == 0 || len(c.Accounts) == 0 {
		return fmt.Errorf("members and accounts must be nonempty")
	}
	members, providers, connectors, accounts, bindings := map[Key]bool{}, map[Key]bool{}, map[Key]bool{}, map[Key]bool{}, map[Key]bool{}
	for i, m := range c.Members {
		path := fmt.Sprintf("members[%d]", i)
		if err := addKey(members, m.Key, path); err != nil {
			return err
		}
		if strings.TrimSpace(m.Name) == "" {
			return fmt.Errorf("%s.name: required", path)
		}
	}
	for i, p := range c.Providers {
		path := fmt.Sprintf("providers[%d]", i)
		if err := addKey(providers, p.Key, path); err != nil {
			return err
		}
		if strings.TrimSpace(p.Name) == "" {
			return fmt.Errorf("%s.name: required", path)
		}
	}
	for i, k := range c.Connectors {
		if err := addKey(connectors, k, fmt.Sprintf("connectors[%d]", i)); err != nil {
			return err
		}
	}
	for i, a := range c.Accounts {
		path := fmt.Sprintf("accounts[%d]", i)
		if err := addKey(accounts, a.Key, path); err != nil {
			return err
		}
		if err := a.definition().ValidateDefinition(); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		owners := map[Key]bool{}
		shares := make([]int, 0, len(a.Owners))
		for j, o := range a.Owners {
			ownerPath := fmt.Sprintf("%s.owners[%d].member", path, j)
			if err := reference(members, o.Member, ownerPath); err != nil {
				return err
			}
			if owners[o.Member] {
				return fmt.Errorf("%s: duplicate owner %q", ownerPath, o.Member)
			}
			owners[o.Member] = true
			shares = append(shares, o.ShareBasisPoints)
		}
		if err := domain.ValidateOwnershipShares(shares); err != nil {
			return fmt.Errorf("%s.owners: %w", path, err)
		}
	}
	byAccount := map[Key][]ProviderBinding{}
	for i, b := range c.Bindings {
		path := fmt.Sprintf("bindings[%d]", i)
		if err := addKey(bindings, b.Key, path); err != nil {
			return err
		}
		if err := reference(accounts, b.Account, path+".account"); err != nil {
			return err
		}
		if err := reference(providers, b.Provider, path+".provider"); err != nil {
			return err
		}
		if err := reference(connectors, b.Connector, path+".connector"); err != nil {
			return err
		}
		if b.ValidFrom.IsZero() {
			return fmt.Errorf("%s.validFrom: required", path)
		}
		if b.ValidTo != nil && !b.ValidTo.After(b.ValidFrom) {
			return fmt.Errorf("%s.validTo: must follow validFrom", path)
		}
		byAccount[b.Account] = append(byAccount[b.Account], b)
	}
	// Iterate accounts in file order, and sort a copied binding slice. Adjacent
	// intervals may touch; open-ended or intersecting intervals may not overlap.
	for _, a := range c.Accounts {
		periods := byAccount[a.Key]
		sort.SliceStable(periods, func(i, j int) bool { return periods[i].ValidFrom.Before(periods[j].ValidFrom) })
		for i := 1; i < len(periods); i++ {
			previous, current := periods[i-1], periods[i]
			if previous.ValidTo == nil || current.ValidFrom.Before(*previous.ValidTo) {
				return fmt.Errorf("account %q: bindings %q and %q overlap", a.Key, previous.Key, current.Key)
			}
		}
	}
	return nil
}
