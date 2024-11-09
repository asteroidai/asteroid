// import { Decision, SupervisionResult, SupervisorChain, Tool, ToolRequestGroup, useGetRunRequestGroups } from '@/types';
// import React, { useEffect, useState } from 'react';
// import { useParams } from 'react-router-dom';
// import Page from './page';
// import ContextDisplay from '@/components/context_display'
// import { UUIDDisplay } from './uuid_display';
// import JsonDisplay from './json_display';
// import { DecisionBadge, ExecutionStatusBadge, StatusBadge, SupervisorBadge, ToolBadge } from './status_badge';
// import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from './ui/accordion';
// import { FileJsonIcon, GitPullRequestIcon, MessagesSquareIcon, PickaxeIcon, PinIcon } from 'lucide-react';
// import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
// import { CreatedAgo } from './created_ago';
// import { Link1Icon } from '@radix-ui/react-icons';


// // TODO allow execution supervision to be passed in
// export default function ExecutionComponent() {
//   const [supervisions, setSupervisions] = useState<ToolRequestGroup[]>();
//   const { executionId } = useParams()
//   const { data, isLoading, isError } = useGetRunRequestGroups(executionId || '');

//   useEffect(() => {
//     if (data?.data) {
//       setSupervisions(data.data);
//     }
//   }, [data]);

//   if (!supervisions) {
//     return (
//       <p>No supervision found</p>
//     )
//   }

//   return (
//     <Page
//       title="Supervision requests & results"
//       subtitle={
//         <span>
//           {/* {supervisions.supervisions.length} supervision chains with {supervisions.supervisions.map(s => s.length).reduce((a, b) => a + b, 0)} requests have been made so far for execution <UUIDDisplay uuid={executionId} /> which is currently in status {` `} */}
//         </span>
//       }
//       icon={<GitPullRequestIcon className="w-6 h-6" />}
//     >
//       {isLoading && (
//         <p>Loading</p>
//       )}
//       {isError && (
//         <p>Error</p>
//       )}
//       <div className="col-span-3 flex flex-col space-y-4">
//         <div>
//           {/* {supervisions.supervisions[0][0].request.tool_requests && (
//             <ToolBadge toolId={supervisions.supervisions[0][0].request.tool_requests[0].tool_id} />
//           )} */}
//         </div>
//         <div>
//           <SupervisionResultsForExecution executionSupervisions={supervisions} />
//         </div>
//       </div >
//     </Page >
//   )
// }

// export function SupervisionsForExecution({ executionId }: { executionId: string }) {
//   const { data, isLoading, isError } = useGetExecutionSupervisions(executionId);

//   if (!data) {
//     return <p>No supervisions found</p>
//   }

//   return (
//     <>
//       {isLoading && <p>Loading</p>}
//       {isError && <p>Error</p>}
//       <SupervisionResultsForExecution executionSupervisions={data.data} />
//     </>
//   )
// }

// export function SupervisionResultsForExecution({ executionSupervisions }: { executionSupervisions: ExecutionSupervisions }) {
//   return (

//     <div className="w-full">
//       {
//         executionSupervisions.supervisions.map((chain, index) => (
//           <div>
//             <SupervisionChainCard chain={chain} chainNumber={index + 1} />
//           </div>
//         ))
//       }
//     </div>
//   )
// }

// function SupervisionResultCard({ result, supervisorId }: { result: SupervisionResult | undefined, supervisorId: string }) {
//   if (!result) {
//     return <p>No result has yet been recorded for this request</p>
//   }

//   return (
//     <Card>
//       <CardHeader>
//         <CardTitle>
//           Supervision Result: <SupervisorBadge supervisorId={supervisorId} /> returned <DecisionBadge decision={result.decision} />
//         </CardTitle>
//         <CardDescription>
//           <CreatedAgo datetime={result.created_at} label="Supervision result occurred" />. ID is <UUIDDisplay uuid={result.id} />
//         </CardDescription>
//       </CardHeader>
//       <CardContent>
//         <p>Reasoning: {result.reasoning != "" ? result.reasoning : "No reasoning given"}</p>
//       </CardContent>
//     </Card>

//   )
// }

// function SupervisionChainCard({ chain, chainNumber }: { chain: SupervisorChain, chainNumber: number }) {
//   return (
//     <div className="w-full space-y-4 bg-muted/50 rounded-md px-4 flex flex-row items-center justify-between gap-4">
//       <div className="flex flex-row items-center justify-center gap-2">
//         <Link1Icon className="w-4 h-4" />
//         <p className="text-xs text-center text-gray-500">{chainNumber}</p>
//       </div>
//       <div className="w-full">
//         {chain.supervisors.map((supervision, index) => (
//           <Accordion type="single" collapsible className="w-full">
//             <AccordionItem value="hub-stats" className="border border-gray-200 rounded-md mb-4">
//               <AccordionTrigger className="w-full p-4 rounded-md cursor-pointer focus:outline-none">
//                 <div className="flex flex-row w-full justify-between">
//                   <div className="flex flex-row gap-2">
//                     <>
//                       <span className="text-sm">Supervision Request #{index + 1} to supervisor</span>
//                       <SupervisorBadge supervisorId={supervision.request.supervisor_id || ''} />
//                       is in status
//                       {supervision.statuses && (
//                         <StatusBadge statuses={supervision.statuses} />
//                       )}
//                       {supervision.result && (
//                         <>
//                           because supervisor decided to
//                           <DecisionBadge decision={supervision.result?.decision} />
//                         </>
//                       )}
//                     </>
//                   </div>


//                   <span>
//                   </span>

//                 </div>
//               </AccordionTrigger>
//               <AccordionContent className="p-4 bg-white rounded-md space-y-4">
//                 <p className="text-xs text-gray-500">
//                   Supervision info for request <UUIDDisplay uuid={supervision.request.id} /> as part of execution <UUIDDisplay uuid={supervision.request.execution_id} />
//                 </p>

//                 <Accordion type="single" collapsible className="w-full">
//                   <AccordionItem value="hub-stats" className="border border-gray-200 rounded-md mb-4">
//                     <AccordionTrigger className="w-full p-4 rounded-md cursor-pointer focus:outline-none">
//                       <div className="flex flex-row gap-4 text-center">
//                         <PickaxeIcon className="w-4 h-4" />
//                         Tool Requests
//                       </div>
//                     </AccordionTrigger>
//                     <AccordionContent>
//                       {supervision.request.tool_requests.map((tool_request, index) => (
//                         <div key={index} className="px-4">
//                           <JsonDisplay json={tool_request} />
//                         </div>
//                       ))}
//                     </AccordionContent>
//                   </AccordionItem>
//                 </Accordion>

//                 <Accordion type="single" collapsible className="w-full">
//                   <AccordionItem value="hub-stats" className="border border-gray-200 rounded-md mb-4">
//                     <AccordionTrigger className="w-full p-4 rounded-md cursor-pointer focus:outline-none">
//                       <div className="flex flex-row gap-4">
//                         <MessagesSquareIcon className="w-4 h-4" />Messages
//                       </div>
//                     </AccordionTrigger>
//                     <AccordionContent>
//                       <ContextDisplay context={supervision.request.task_state} />
//                     </AccordionContent>
//                   </AccordionItem>
//                 </Accordion>

//                 <Accordion type="single" collapsible className="w-full">
//                   <AccordionItem value="hub-stats" className="border border-gray-200 rounded-md mb-4">
//                     <AccordionTrigger className="w-full p-4 rounded-md cursor-pointer focus:outline-none">
//                       <div className="flex flex-row gap-4 text-center">
//                         <FileJsonIcon className="w-4 h-4" />
//                         Full Task State JSON
//                       </div>
//                     </AccordionTrigger>
//                     <AccordionContent>
//                       <JsonDisplay json={supervision.request} />
//                     </AccordionContent>
//                   </AccordionItem>
//                 </Accordion>

//                 <SupervisionResultCard result={supervision.result} supervisorId={supervision.request.supervisor_id || ''} />
//               </AccordionContent>
//             </AccordionItem>
//           </Accordion>
//         ))}
//       </div>
//     </div>
//   )
// }
