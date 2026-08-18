import { getSecurityMetadata } from "@/utils.js";
import {ZUserDevice} from "@glimpse/zod"
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const deviceContract = c.router(
  {
    registerDevice: {
      summary: "Register a user device",
      path: "/register",
      method: "POST",
      description: "Registers a device for push notifications against the calling user",
      body: z.object({
        deviceToken: z.string().min(1),
        platform: z.string().min(1),
      }),
      responses: {
        201: ZUserDevice,
      },
      metadata: metadata,
    }
  }
)
