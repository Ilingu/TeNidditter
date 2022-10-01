import crypto from "crypto";

import type { FunctionJob } from "$lib/types/interfaces";
import { error } from "@sveltejs/kit";

import { ENCRYPT_KEY } from "$env/static/private";

export const GET: import("./$types").RequestHandler = ({ url }) => {
	const todec = url.searchParams.get("dec");

	if (!todec) throw error(400);
	const toDec = decodeURIComponent(todec);

	const { success, data: decoded } = decryptDatas(Buffer.from(toDec, "base64"));
	if (!success || !decoded) throw error(400);

	return new Response(String(encodeURIComponent(decoded.toString("utf8"))));
};

/**
 * Decrypt data
 * @param {Buffer} encryptedStr
 */
const decryptDatas = (encryptedStr: Buffer): FunctionJob<Buffer> => {
	try {
		const iv = encryptedStr.subarray(3, 3 + 12);
		const ciphertext = encryptedStr.subarray(3 + 12, encryptedStr.length - 16);
		const authTag = encryptedStr.subarray(encryptedStr.length - 16);
		const decipher = crypto.createDecipheriv("aes-256-gcm", Buffer.from(ENCRYPT_KEY), iv);
		decipher.setAuthTag(authTag);
		const decryptedCookie = Buffer.concat([decipher.update(ciphertext), decipher.final()]);
		return { success: true, data: decryptedCookie }; // .toString("utf8")
	} catch (err) {
		console.log(err);
		return { success: false, error: JSON.stringify(err) };
	}
};
