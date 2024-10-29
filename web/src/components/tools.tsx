import { Tool, useGetProject, useGetTools } from "@/types";
import React, { useEffect, useState } from "react";
import Page from "./page";
import { CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "./ui/card";
import { Card } from "@radix-ui/themes";
import { Link } from "react-router-dom";
import { ToolsList } from "@/components/tools_list";
import { useProject } from '@/contexts/project_context';

export default function Tools() {
  const [tools, setTools] = useState<Tool[]>([]);
  const { selectedProject } = useProject();
  const { data: projectData, isLoading: projectLoading, error: projectError } = useGetProject(selectedProject || '');

  const { data, isLoading, error } = useGetTools(
    { projectId: selectedProject! },
  );

  useEffect(() => {
    if (data?.data) {
      setTools(data.data);
    }
  }, [data]);


  if (!selectedProject) {
    return <div>Please select a project first</div>;
  }

  if (isLoading) return <Page title="Tools">Loading...</Page>;
  if (error) return <Page title="Tools">Error: {error.message}</Page>;

  return (
    <Page title={`Tools for project ${projectData?.data?.name}`} subtitle={`${tools.length} tools found for project ${projectData?.data?.id}`}>
      {tools.length > 0 && <ToolsList tools={tools} variant="card" />}
      {tools.length === 0 && <p className="text-sm text-gray-500">No tools found for this project. When your agent registers a tool, it will appear here.</p>}

    </Page>
  );
}
