export function money(value: number): string {
  return new Intl.NumberFormat('en-GB', { style: 'currency', currency: 'GBP' }).format(value);
}

export function percent(value: number, total: number): string {
  if (total === 0) {
    return '0%';
  }
  return `${Math.round((value / total) * 100)}%`;
}

