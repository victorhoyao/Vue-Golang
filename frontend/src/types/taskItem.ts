export interface TaskItem {
  id?: number;
  taskId: string;
  taskName: string;
  supplier: string;
  userTypes: string;
  priceReal: number;
  priceProtocol: number;
  priceManual: number;
  selectionMode: string;
  status: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface TaskItemFormData extends Omit<TaskItem, 'id' | 'createdAt' | 'updatedAt'> {
  // Any specific form-only fields or modifications
}

export interface TaskItemRes {
  list: TaskItem[];
  total: number;
} 