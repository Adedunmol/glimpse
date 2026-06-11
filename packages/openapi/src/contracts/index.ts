import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { uploadContract } from "./upload.js";

const c = initContract();

export const apiContract = c.router({
  Health: healthContract,
  Upload: uploadContract,
});