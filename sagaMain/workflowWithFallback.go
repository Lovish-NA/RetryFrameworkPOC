package main

struct{
	int step

}

function syncWorkflowWithFallback(activity,fallbackActivity){
	err := workflow.execute().;
	if(err != nil){
		fallback()
	}
	stack.push(fallbackActivity);

}
