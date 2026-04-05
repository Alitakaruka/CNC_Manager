export function getCncUniqueKey(cnc) {
  if (!cnc || typeof cnc !== 'object') return ''
  return cnc.uniqueKey || cnc.UniqueKey || ''
}

export function mergeCncIntoList(prev, cnc) {
  const key = getCncUniqueKey(cnc)
  if (!key) return prev
  const index = prev.findIndex(c => getCncUniqueKey(c) === key)
  if (index !== -1) {
    const next = [...prev]
    next[index] = { ...next[index], ...cnc }
    return next
  }
  return [...prev, cnc]
}

/** Полный снимок (массив из нескольких станков) заменяет state; один объект или [один] — мерж по uniqueKey. */
export function applyPrintersWsPayload(prev, data) {
  if (Array.isArray(data)) {
    if (data.length === 0) return prev
    if (data.length === 1) return mergeCncIntoList(prev, data[0])
    return data
  }
  if (data && typeof data === 'object' && getCncUniqueKey(data)) {
    return mergeCncIntoList(prev, data)
  }
  return prev
}
