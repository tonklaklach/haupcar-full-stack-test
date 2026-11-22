const platePattern = /^([0-9]?[ก-ฮ]{2})-?([0-9]{1,4})$/u;

export const registrationError =
  "รูปแบบทะเบียนรถไม่ถูกต้อง เช่น 4ขห-798 หรือ ขก-1234";

export function formatRegistration(value: string): string | null {
  const parts = platePattern.exec(value.replace(/\p{White_Space}/gu, ""));
  return parts ? `${parts[1]}-${parts[2]}` : null;
}

export function registrationSearchKey(value: string) {
  return value.toLowerCase().replace(/[\p{White_Space}-]/gu, "");
}
