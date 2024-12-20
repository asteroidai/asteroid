import { Tool, useGetProject, useGetProjectTools } from "@/types";
import React, { useEffect, useState } from "react";
import Page from "./util/page";
import { ToolsList } from "@/components/tools_list";
import { useProject } from '@/contexts/project_context';
import { UUIDDisplay } from "./util/uuid_display";
import LoadingSpinner from "./util/loading";
import { PickaxeIcon } from "lucide-react";
import { ToolCard } from "./tool_card";

export default function Tools() {
  const [tools, setTools] = useState<Tool[]>([]);
  const { selectedProject } = useProject();
  const { data: projectData, isLoading: projectLoading, error: projectError } = useGetProject(selectedProject || '');

  const { data, isLoading, error } = useGetProjectTools(
    selectedProject || '',
    {
      query: {
        enabled: !!selectedProject,
        refetchInterval: 1000,
      },
    }
  );

  useEffect(() => {
    if (data?.data) {
      setTools(data.data);
    } else {
      setTools([]);
    }
  }, [data, selectedProject]);


  if (!selectedProject) {
    return <div>Please select a project first</div>;
  }

  return (
    <Page
      title={`Tools`}
      subtitle={<span>{tools.length > 0 ? `${tools.length} tool${tools.length === 1 ? '' : 's'}` : 'No tools'} found for project {projectData?.data?.name ?? ''} <UUIDDisplay uuid={projectData?.data?.id ?? ''} /></span>}
      icon={<PickaxeIcon className="w-6 h-6" />}
      cols={2}
    >
      {tools.length === 0 && (
        <div className="col-span-2">
          <p className="text-sm text-gray-500">No tools found for this project. When your agent registers a tool, it will appear here.</p>
        </div>
      )}
      {isLoading && (
        <LoadingSpinner />
      )}
      {error && (
        <div>Error loading tools: {error.message}</div>
      )}
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">

        {(tools as Tool[]).map((tool) => (
          <div className="col-span-1">
            <ToolCard key={tool.id} tool={tool} />
          </div>
        ))}
      </div>
    </Page>
  );
}
