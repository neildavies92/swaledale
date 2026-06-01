-- name: ListMembers :many
SELECT id, name FROM members ORDER BY id;

-- name: UpdateBudgetItem :one
UPDATE budget_items
SET label = $3, amount_pence = $4
WHERE member_id = $1 AND id = $2
RETURNING id, member_id, category_id, label, amount_pence;

-- name: UpdateJointAccountItem :one
UPDATE joint_account_items
SET label = $2, amount_pence = $3
WHERE id = $1
RETURNING id, label, amount_pence;

-- name: UpdateGoal :one
UPDATE household_goals
SET name = $2, target_pence = $3, current_pence = $4, notes = $5
WHERE id = $1
RETURNING id, name, target_pence, current_pence, notes;

