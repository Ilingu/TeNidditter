import { IsEmptyString, IsValidJSON } from "$lib/utils";
import type { RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ url }) => {
	const rawRecoveryCodes = decodeURIComponent(url.searchParams.get("codes") ?? "");
	if (IsEmptyString(rawRecoveryCodes) || !IsValidJSON(rawRecoveryCodes))
		return new Response("no codes provided", { status: 400 });

	const recoveryCodes: string[] = JSON.parse(rawRecoveryCodes);

	return new Response(
		` <!DOCTYPE html>
      <html lang="en">
      <head>
      	<meta charset="UTF-8">
      	<meta http-equiv="X-UA-Compatible" content="IE=edge">
      	<meta name="viewport" content="width=device-width, initial-scale=1.0">
      	<title>🔒 Twiniditter Recovery Codes...</title>
        <style>
          * {
              font-family: monospace;
              box-sizing: border-box;
              margin: 0;
              padding: 0;
            }
            body {
              min-height: 100vh;
              display: flex;
              flex-direction: column;
              justify-content: center;
              align-items: center;
              background: #fff;
              color: #000;
            }
            body > h1 {
              margin-bottom: 10px;
            }
            body > p {
              margin-top: 10px;
            }
        </style>
      </head>
      <body>
            <h1>🔒 Twiniditter Recovery codes</h1>
            <ul>
              ${recoveryCodes?.map((code) => `<li>${code}</li>`).join("")}
            </ul>
            <p>Twiniditter two-factor authentication account recovery codes.</p>

            <script>
                window.addEventListener('afterprint', (event) => {
                    const win = window.open("/", "_self");
                    win.close();
                });

	              window.onload = print;
            </script>
      </body>
      </html>`,
		{
			status: 200,
			headers: { "Content-Type": "text/html" }
		}
	);
};
