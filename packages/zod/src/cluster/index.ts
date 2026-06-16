import z from "zod";

export const ZCluster = z.object({
  id: z.string().uuid(),
  uploadId: z.string(),
  label: z.string(),
  thumbnailPhototId: z.string().uuid(),
  createdAt: z.string(),
  updatedAt: z.string(),
})