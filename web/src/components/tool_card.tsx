import { SupervisorChain, Tool, useGetToolSupervisorChains } from "@/types";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "./ui/card";
import React, { useEffect, useState } from "react";
import { ArrowRightIcon } from "@radix-ui/react-icons";
import { SupervisorBadge, SupervisorTypeBadge, ToolBadge } from "./util/status_badge";
import { ToolAttributes } from "./tool_attributes";
import { UUIDDisplay } from "./util/uuid_display";
import hljs from 'highlight.js';
import 'highlight.js/styles/github-dark.css';

type ToolCardProps = {
  tool: Tool;
  runId?: string;
};

export function ToolCard({ tool, runId }: ToolCardProps) {
  useEffect(() => {
    if (tool.code) {
      document.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightElement(block as HTMLElement);
      });
    }
  }, [tool.code]);

  return (
    <Card className="flex flex-col text-muted-foreground">
      <CardHeader className="py-2">
        <CardTitle className="py-4 flex flex-row justify-between">
          <ToolBadge toolId={tool.id || ''} />
          <UUIDDisplay uuid={tool.id || ''} className="text-xs font-normal" />
        </CardTitle>
        <CardDescription>
          <p className="text-sm font-semibold">Description</p>
          {/* <textarea className="text-xs bg-muted p-2 rounded w-full resize-none" value={JSON.stringify(tool.description, null, 2)} readOnly rows={10} disabled /> */}
          {tool.description && <p className="text-sm">{tool.description}</p>}
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-2">
        {tool.code && (
          <div className="">
            <p className="text-sm font-semibold mb-2">Code</p>
            <pre className="text-xs rounded w-full overflow-scroll max-h-96">
              <code className="rounded-lg language-python">{tool.code}</code>
            </pre>
          </div>
        )}
        <p className="text-sm font-semibold">Attributes</p>
        {tool.attributes && <ToolAttributes attributes={tool.attributes || ''} ignoredAttributes={tool.ignored_attributes || []} />}
        {runId && tool.id && <RunToolSupervisors runId={runId} toolId={tool.id} />}
      </CardContent>
    </Card>
  );
}

function RunToolSupervisors({ runId, toolId }: { runId: string, toolId: string }) {
  const [supervisorChain, setSupervisorChain] = useState<SupervisorChain[]>([]);
  const { data, isLoading, error } = useGetToolSupervisorChains(toolId);

  useEffect(() => {
    if (data) {
      setSupervisorChain(data.data);
    }
  }, [data]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <div className="flex flex-col gap-2">
      <p className="text-sm font-semibold">Supervisors configured for this tool</p>
      {supervisorChain.map((chain, index) => (
        <div className="flex flex-row overflow-scroll gap-2 bg-muted p-2 rounded-md items-center" key={chain.chain_id}>
          <p className="text-sm font-semibold text-muted-foreground">Chain {index + 1}</p>
          {chain.supervisors.map((supervisor, index) => (
            <div key={supervisor.id} className="flex flex-row gap-2 items-center">
              {index > 0 && <ArrowRightIcon className="w-4 h-4" />}
              {supervisor.id && <SupervisorBadge supervisorId={supervisor.id} />}
            </div>
          ))}
        </div>
      ))}
    </div>
  );
}
