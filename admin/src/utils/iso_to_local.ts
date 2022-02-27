const localeOptions = {
  year: "numeric",
  month: "long",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
  second: "numeric",
  hour12: false,
} as Intl.DateTimeFormatOptions;

export function IsoDateToLocalStr(isoDate: Date): string {
  const date = new Date(isoDate);
  return date.toLocaleString(navigator.language, localeOptions);
}
