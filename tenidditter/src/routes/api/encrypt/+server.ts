import crypto from "crypto";

import type { FunctionJob } from "$lib/types/interfaces";
import { error } from "@sveltejs/kit";

import { ENCRYPT_KEY } from "$env/static/private";

export const GET: import("./$types").RequestHandler = ({ url }) => {
	const toenc = url.searchParams.get("enc");

	if (!toenc) throw error(400);
	const toEnc = decodeURIComponent(toenc);

	const { success, data: encoded } = encryptDatas(Buffer.from(toEnc));
	if (!success || !encoded) throw error(400);

	return new Response(String(encodeURIComponent(encoded.toString("base64"))));
};

/**
 * Encrypt data
 * @param {Buffer} str
 */
const encryptDatas = (str: Buffer): FunctionJob<Buffer> => {
	try {
		const iv = crypto.randomBytes(12);
		const cipher = crypto.createCipheriv("aes-256-gcm", Buffer.from(ENCRYPT_KEY), iv);
		const encryptedStr = Buffer.concat([
			Buffer.from("v10"),
			iv,
			cipher.update(str),
			cipher.final(),
			cipher.getAuthTag()
		]);
		return { success: true, data: encryptedStr };
	} catch (err) {
		console.error(err);
		return { success: false, error: JSON.stringify(err) };
	}
};
