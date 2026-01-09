// Shared Types for Nexus Ops Frontend

export interface AppError {
    code: string;
    message: string;
}

export interface MetricsDTO {
    activeGoroutines: number;
    tasksCompleted: number;
    errors: number;
}

export interface Task {
    id: string;
    name: string;
    type: string;
    startedAt: string;
    status: string;
}
