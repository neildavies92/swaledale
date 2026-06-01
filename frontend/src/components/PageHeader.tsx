export function PageHeader({ title, eyebrow }: { title: string; eyebrow?: string }) {
  return (
    <header className="mb-6">
      {eyebrow ? <p className="text-sm font-medium text-moss">{eyebrow}</p> : null}
      <h2 className="text-2xl font-semibold tracking-normal sm:text-3xl">{title}</h2>
    </header>
  );
}

