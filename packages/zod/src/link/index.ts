import z from "zod";

export const ZLink = z.object({
  id: z.string().uuid(),
  clusterId: z.string().uuid(),
  token: z.string(),
  isPasswordProtected: z.boolean(),
  passwordHash: z.string(),
  expiresAt: z.string(),
  isActive: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
})