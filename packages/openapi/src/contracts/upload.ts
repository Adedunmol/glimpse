import { getSecurityMetadata } from "@/utils.js";
import {schemaWithPagination, ZPresignedUrls, ZUpload} from "@glimpse/zod";
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const uploadContract = c.router(
    {
        getUploads: {
            summary: "Get all uploads",
            path: "/uploads",
            method: "GET",
            description: "Get all uploads",
            query: z.object({
                page: z.number().min(1).optional(),
                limit: z.number().min(1).max(100).optional(),
                sort: z.enum(["created_at", "updated_at", "name"]).optional(),
                order: z.enum(["asc", "desc"]).optional(),
                search: z.string().min(1).optional(),
                status: ZUpload.shape.status.optional()
            }),
            responses: {
                200: schemaWithPagination(ZUpload)
            },
            metadata: metadata
        },
        
        createUpload: {
        summary: "Create a new upload",
        path: "/uploads",
        method: "POST",
        description: "Create a new upload",
        body: ZUpload.pick({
            name: true,
            expiresAt: true
        })
            .partial(),
        responses: {
            201: ZUpload,
        },
        metadata: metadata,
        },

        getUploadById: {
        summary: "Get upload by ID",
        path: "/uploads/:id",
        method: "GET",
        description: "Get upload by ID",
        responses: {
            200: ZUpload,
        },
        metadata: metadata,
        },

        updateUpload: {
        summary: "Update upload",
        path: "/uploads/:id",
        method: "PATCH",
        description: "Update upload",
        body: ZUpload.pick({
            name: true,
            expiresAt: true
        }).partial(),
        responses: {
            200: ZUpload,
        },
        metadata: metadata,
        },

        deleteUpload: {
        summary: "Delete upload",
        path: "/uploads/:id",
        method: "DELETE",
        description: "Delete upload",
        responses: {
            204: z.void(),
        },
        metadata: metadata,
        },

        getPresignedUrls: {
            summary: "Get presigned urls",
            path: "/uploads/:id/photos",
            method: "POST",
            description: "This endpoint takes in a list of the images to be uploaded with the upload id and generates a presigned url for each of the photo in the payload",
            body: z.object({
                files: z.array(z.object({
                    name: z.string()
                })),
            }),
            responses: {
                200: ZPresignedUrls,
            },
            metadata: metadata,
        },

        completeUpload: {
        summary: "Complete upload",
        path: "/uploads/:id/complete",
        method: "POST",
        description: "This endpoint takes in a list of pictures that have been successfully uploaded to S3 by the client and creates records for them in the database and also kickstart the processing flow",
        body: z.object({
            files: z.array(z.object({
                name: z.string()
            })),
        }),
        responses: {
            204: z.void(),
        },
        metadata: metadata,
        },
    },
    {
        pathPrefix: "/v1"
    }
)