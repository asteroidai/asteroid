import { Supervisor, Tool, useGetRunToolSupervisors } from "@/types";
import { Link } from "react-router-dom";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "./ui/card";
import { Badge } from "./ui/badge";
import { Button } from "./ui/button";
import React, { useEffect, useState } from "react";
import { ArrowRightIcon } from "@radix-ui/react-icons";

type ToolCardProps = {
  tool: Tool;
  runId?: string;
};

export function ToolCard({ tool, runId }: ToolCardProps) {
  return (
    <Card className="flex flex-col ">
      <CardHeader>
        <CardTitle>{tool.name}</CardTitle>
        <CardDescription>
          {tool.description}
          <div>Tool ID: {tool.id}</div>
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow flex flex-col gap-4">

        {tool.created_at}
        {tool.attributes && <pre className="text-xs mt-1 bg-muted p-2 rounded">{JSON.stringify(tool.attributes, null, 2)}</pre>}
        {runId && <RunToolSupervisors runId={runId} toolId={tool.id} />}
      </CardContent>

      <CardFooter className="flex justify-end">
        <Link to={`/tools/${tool.id}`} key={tool.id}>
          <Button variant="outline">View Tool <ArrowRightIcon className="w-4 h-4" /></Button>
        </Link>
      </CardFooter>

    </Card>
  );
}

function RunToolSupervisors({ runId, toolId }: { runId: string, toolId: string }) {
  const [supervisors, setSupervisors] = useState<Supervisor[]>([]);
  const { data, isLoading, error } = useGetRunToolSupervisors(runId, toolId);

  useEffect(() => {
    if (data) {
      setSupervisors(data.data);
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
      <p className="text-sm font-semibold">Supervisors</p>
      {supervisors.map((supervisor) => (
        <Card key={supervisor.id}>
          <CardHeader>
            <CardTitle className="flex flex-row justify-between">
              {supervisor.id}
              <Badge variant="secondary" className="text-xs">{supervisor.type}</Badge>
            </CardTitle>
            <CardDescription>
              {supervisor.description}
            </CardDescription>


          </CardHeader>
        </Card>
      ))}
    </div>
  );
}