export interface BarcodeScannerOption {
    scannerId: number;
    status: string;
}

export interface ScannerResponse {
    messageType: "response" | "barcode" | "image" | "attached" | "detached" | "scanner_list";
    scannerId: number;
    barcode?: string;
    image?: string;
    scanners?: BarcodeScannerOption[];
    status: boolean;
}