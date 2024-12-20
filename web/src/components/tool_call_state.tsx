import { useGetToolCallState } from '@/types';
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@radix-ui/react-accordion';
import { Card } from '@radix-ui/themes';
import { DrillIcon } from 'lucide-react';
import * as React from 'react';
import { CardContent } from './ui/card';
import RunExecutionViewer from './chain_execution_viewer';
import ToolCallCard from './tool_call_card';
import JsonDisplay from './util/json_display';
import { useEffect } from 'react';
import { useState } from 'react';

interface ToolCallStateProps {
  toolCallId: string | undefined
}

export function ToolCallState({ toolCallId }: ToolCallStateProps) {
  const { data: toolCallStateData } = useGetToolCallState(toolCallId || '');
  const [expanded, setExpanded] = useState<boolean>(false);

  useEffect(() => {
    if (toolCallStateData?.data) {
      setExpanded(true);
    }
  }, [toolCallStateData]);

  return (
    <>
      {/* <Button onClick={() => setExpanded(!expanded)}>
        {expanded ? "Hide" : "Show"} Details
      </Button> */}
      <Accordion
        type="single"
        collapsible
        className="w-full"
        value={expanded ? "messages" : undefined}
        onValueChange={(value) => setExpanded(!!value)}
      >
        <AccordionItem value="messages" className="border border-gray-200 rounded-md">
          <AccordionTrigger className="w-full p-4 rounded-md cursor-pointer focus:outline-none">
            <div className="flex flex-row gap-4 items-center">
              <DrillIcon className="w-4 h-4" />
              Tool Call Information
            </div>
          </AccordionTrigger>
          <AccordionContent className="">
            <Card className="border-none">
              <CardContent className="space-y-4">
                {toolCallId ? (
                  <>
                    <ToolCallCard status={toolCallStateData?.data?.status as "completed" | "failed" | "pending"} toolCall={toolCallStateData?.data?.toolcall} />
                    <RunExecutionViewer runExecution={toolCallStateData?.data} />
                    <JsonDisplay json={toolCallStateData?.data} />
                  </>
                ) : (
                  <div className="text-sm text-muted-foreground">No tool call selected. Click on a tool call in the conversation to view details.</div>
                )}
              </CardContent>
            </Card>
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </>
  );
}
