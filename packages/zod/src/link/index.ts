import z from "zod";

export const ZLink = z.object({
  id: z.string().uuid(),
  clusterId: z.string().uuid(),
  token: z.string(),
  isPasswordProtected: z.boolean(),
  passwordHash: z.string().nullable(),
  expiresAt: z.string(),
  isActive: z.boolean(),
  createdAt: z.string(),
  updatedAt: z.string(),
})