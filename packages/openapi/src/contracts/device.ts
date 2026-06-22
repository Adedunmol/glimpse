import { getSecurityMetadata } from "@/utils.js";
import {ZUserDevice} from "@glimpse/zod"
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const deviceContract = c.router(
  {
    createDevice: {
      summary: "create a new user device",
      path: "/device/create",
      method: "POST",
      description:"create a new user device",
      body: ZUserDevice.pick({
        pushToken: true,
        platform: true
      }).partial(),
      responses: {
        201: ZUserDevice,
      },
      metadata: metadata,
    }
  }
)