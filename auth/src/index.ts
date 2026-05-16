import { serve } from "@hono/node-server";
import { bootstrap } from "~/composition/root";

const { app } = bootstrap();

serve(app);
