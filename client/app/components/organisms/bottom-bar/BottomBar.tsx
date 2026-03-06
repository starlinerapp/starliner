import React, { useEffect, useRef, useState } from "react";
import type {
  ResponseDatabaseDeployment,
  ResponseGitDeployment,
  ResponseImageDeployment,
  ResponseIngressDeployment,
} from "~/server/api/client/generated";

type Deployment =
  | ResponseGitDeployment
  | ResponseImageDeployment
  | ResponseIngressDeployment
  | ResponseDatabaseDeployment;

interface BottomBarProps {
  deployment?: Deployment;
}

export default function BottomBar({ deployment }: BottomBarProps) {
  const lastDeploymentRef = useRef<Deployment | undefined>(deployment);

  const containerRef = useRef<HTMLDivElement>(null);
  const titleRef = useRef<HTMLSpanElement>(null);

  const [underline, setUnderline] = useState({ left: 0, width: 0 });

  useEffect(() => {
    if (deployment) {
      lastDeploymentRef.current = deployment;
    }
  }, [deployment]);

  const selectedDeployment = deployment ?? lastDeploymentRef.current;

  useEffect(() => {
    if (!containerRef.current || !titleRef.current) return;

    const rect = titleRef.current.getBoundingClientRect();
    const containerRect = containerRef.current.getBoundingClientRect();

    setUnderline({
      left: rect.left - containerRect.left,
      width: rect.width,
    });
  }, [selectedDeployment]);

  return (
    <>
      <div className="bg-violet-1">
        <div
          ref={containerRef}
          className="border-mauve-6 text-mauve-11 relative flex w-full gap-4 border-b px-2 pt-2 pb-1 text-sm"
        >
          <div className="relative z-10 px-2 py-1.5">
            <span
              ref={titleRef}
              className="text-violet-11 truncate pb-2 font-semibold"
            >
              {selectedDeployment?.serviceName ?? "Logs"}
            </span>
          </div>

          <div
            className="bg-violet-11 absolute bottom-0 h-[3px] rounded-md"
            style={{
              left: underline.left,
              width: underline.width,
            }}
          />
        </div>
      </div>

      <div className="h-full p-4">
        {!selectedDeployment ? (
          <p className="text-mauve-11 italic">
            No deployment selected. Select one to view logs.
          </p>
        ) : (
          <p className="text-mauve-12 h-full w-full font-mono text-sm">
            2010-09-16T15:13:46.677020+00:00 app[web.1]: Processing
            PostController#list (for 208.39.138.12 at 2010-09-16 15:13:46) [GET]
            2010-09-16T15:13:46.677023+00:00 app[web.1]: Rendering template
            within layouts/application 2010-09-16T15:13:46.677902+00:00
            app[web.1]: Rendering post/list 2010-09-16T15:13:46.678990+00:00
            app[web.1]: Rendered includes/_header (0.1ms)
            2010-09-16T15:13:46.698234+00:00 app[web.1]: Completed in 74ms
            (View: 31, DB: 40) | 200 OK
            [http://example-app-1234567890ab.heroku.com/]
            2010-09-16T15:13:46.723498+00:00 heroku[router]: at=info method=GET
            path="/posts" host=example-app-1234567890ab.herokuapp.com"
            fwd="204.204.204.204" dyno=web.1 connect=1ms service=18ms status=200
            bytes=975 2010-09-16T15:13:47.893472+00:00 app[worker.1]: 2 jobs
            processed at 16.6761 j/s, 0 failed ...
          </p>
        )}
      </div>
    </>
  );
}
