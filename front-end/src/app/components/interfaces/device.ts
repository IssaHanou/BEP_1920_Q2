export interface Device {
  id: string;
  description: string;
  input: string[];
  output: string[];
  status: string[];
  connection: boolean;
}
