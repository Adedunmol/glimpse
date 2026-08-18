import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { uploadContract } from "./upload.js";
import { clerkWebHookContract } from "./clerk.js";
import { clusterContract } from "./cluster.js";
import { linkContract } from "./link.js";
import { deviceContract } from "./device.js";


const c = initContract();

export const apiContract = c.router({
  Health: healthContract,
  Upload: uploadContract,
  Cluster: clusterContract,
  Link: linkContract,
  Device: deviceContract,
  ClerkWebhook: clerkWebHookContract
});
