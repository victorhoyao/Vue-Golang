import axios from 'axios';
import type { TaskItem, TaskItemFormData, TaskItemRes } from '../types/taskItem';

export interface TaskItemCreateData {
  taskId: string;
  taskName: string;
  supplier: string;
  userTypes: string;
  priceReal: number;
  priceProtocol: number;
  priceManual: number;
  selectionMode: string;
  status: string;
}

export interface TaskItemUpdateData {
  taskId?: string;
  taskName: string;
  supplier: string;
  userTypes: string;
  priceReal: number;
  priceProtocol: number;
  priceManual: number;
  selectionMode: string;
  status: string;
}

export interface BackendResponse<T> {
  code: number;
  message: string;
  data: T;
}

// Get all task items
export function getTaskItems() {
  return axios.get<BackendResponse<{ taskItems: TaskItem[] }>>('/TaskItem/list');
}

// Add a new task item
export function addTaskItem(data: TaskItemCreateData) {
  return axios.post<BackendResponse<{ taskItem: TaskItem }>>('/TaskItem/add', data);
}

// Update an existing task item
export function updateTaskItem(id: number, data: TaskItemUpdateData) {
  return axios.put<BackendResponse<{ taskItem: TaskItem }>>(`/TaskItem/update/${id}`, data);
}

// Delete a task item
export function deleteTaskItem(id: number) {
  return axios.delete<BackendResponse<Record<string, never>>>(`/TaskItem/delete/${id}`);
} 