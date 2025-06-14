import axios from 'axios';

export interface BackendResponse<T> {
  code: number;
  msg: string;
  data: T;
}

export interface Supplier {
  id?: number;
  supplierId: string;
  supplierName: string;
  contact?: string;
  status: string;
  tasksAssigned?: number;
}

export interface SupplierCreateData {
  supplierId: string;
  supplierName: string;
  contact?: string;
  status?: string;
  tasksAssigned?: number;
}

export interface SupplierUpdateData {
  supplierId?: string;
  supplierName?: string;
  contact?: string;
  status?: string;
  tasksAssigned?: number;
}

export function addSupplier(data: SupplierCreateData) {
  return axios.post<BackendResponse<Supplier>>('/Supplier/add', data);
}

export function getSuppliers() {
  return axios.get<BackendResponse<{ suppliers: Supplier[] }>>('/Supplier/list');
}

export function updateSupplier(id: number, data: SupplierUpdateData) {
  return axios.put<BackendResponse<Supplier>>(`/Supplier/update/${id}`, data);
}

export function deleteSupplier(id: number) {
  return axios.delete<BackendResponse<any>>(`/Supplier/delete/${id}`);
} 