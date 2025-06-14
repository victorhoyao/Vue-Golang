import axios from 'axios';

export interface TaskDistribution {
  id: number;
  batchName: string;
  taskItemId: number;
  taskItem: {
    id: number;
    taskId: string;
    taskName: string;
    supplier: string;
  };
  totalTasks: number;
  realMachineTasks: number;
  protocolTasks: number;
  manualTasks: number;
  distributedTasks: number;
  remainingTasks: number;
  status: string;
  createdBy: number;
  createdByUser: {
    id: number;
    userName: string;
  };
  createdAt: string;
  updatedAt: string;
}

export interface TaskDistributionCreateData {
  batchName: string;
  taskItemId: number;
  totalTasks: number;
  realMachineTasks: number;
  protocolTasks: number;
  manualTasks: number;
}

export interface TaskDistributionUpdateData {
  batchName?: string;
  totalTasks?: number;
  realMachineTasks: number;
  protocolTasks: number;
  manualTasks: number;
  status?: string;
}

export interface BackendResponse<T> {
  code: number;
  msg: string;
  data: T;
}

export interface PaginationResponse<T> {
  distributions: T[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
  };
}

// Create a new task distribution
export function createTaskDistribution(data: TaskDistributionCreateData) {
  return axios.post<BackendResponse<{ distribution: TaskDistribution }>>('/TaskDistribution/create', data);
}

// Get all task distributions with pagination
export function getTaskDistributions(page = 1, pageSize = 10) {
  return axios.get<BackendResponse<{
    distributions: TaskDistribution[];
    pagination: {
      page: number;
      pageSize: number;
      total: number;
    };
  }>>('/TaskDistribution/list', {
    params: { page, pageSize }
  });
}

// Get a specific task distribution by ID
export function getTaskDistribution(id: number) {
  return axios.get<BackendResponse<{ distribution: TaskDistribution }>>(`/TaskDistribution/get/${id}`);
}

// Update an existing task distribution
export function updateTaskDistribution(id: number, data: TaskDistributionUpdateData) {
  return axios.put<BackendResponse<{ distribution: TaskDistribution }>>(`/TaskDistribution/update/${id}`, data);
}

// Delete a task distribution
export function deleteTaskDistribution(id: number) {
  return axios.delete<BackendResponse<Record<string, never>>>(`/TaskDistribution/delete/${id}`);
}

// Activate a task distribution
export function activateTaskDistribution(id: number) {
  return axios.post<BackendResponse<{ distribution: TaskDistribution }>>(`/TaskDistribution/activate/${id}`);
}

// Get task distributions by task item ID
export function getTaskDistributionsByTaskItem(taskItemId: number) {
  return axios.get<BackendResponse<{ distributions: TaskDistribution[] }>>(`/TaskDistribution/by-task-item/${taskItemId}`);
}

// Get distribution summary statistics
export function getDistributionSummary() {
  return axios.get<BackendResponse<{ summary: any }>>('/TaskDistribution/summary');
} 