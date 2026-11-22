export type Car = {
  id: number;
  registration: string;
  brand: string;
  model: string;
  notes: string;
  created_at: string;
  updated_at: string;
};
export type CarInput = Pick<Car, "registration" | "brand" | "model" | "notes">;
async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`/api/cars${path}`, {
    ...init,
    headers: { "Content-Type": "application/json" },
    cache: "no-store",
  });
  if (!response.ok) {
    const body = await response.json().catch(() => null);
    throw new Error(body?.message ?? "เชื่อมต่อระบบไม่สำเร็จ กรุณาลองใหม่");
  }
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}
export const carsAPI = {
  list: () => request<Car[]>(""),
  save: (input: CarInput, id?: number) =>
    request<Car>(id ? `/${id}` : "", {
      method: id ? "PUT" : "POST",
      body: JSON.stringify(input),
    }),
  remove: (id: number) => request<void>(`/${id}`, { method: "DELETE" }),
};
