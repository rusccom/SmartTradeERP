export function rateChanges(current, initial = {}) {
  const changes = {};
  for (const country of new Set([...Object.keys(current), ...Object.keys(initial)])) {
    if (current[country] !== initial[country]) changes[country] = current[country] ?? null;
  }
  return changes;
}
