import { useEffect, useRef } from "react";
import type { ScannerResponse } from "./BarcodeScanner";
import { toast } from "react-toastify";

type Props = {
    barcodeData?: string;
    onBarcodeScan?: (scannerId: number | undefined, barcode: string) => void;
};

export default function BarcodeListener({
    barcodeData,
    onBarcodeScan,
}: Props) {
    const mountedRef = useRef<boolean>(false);
    const keyPresses = useRef<string>("");
    const intervalRef = useRef<number>(undefined);

    const processKeypresses = (scannerId: number | undefined, keyPresses: string) => {
        const matches = keyPresses.match(/^(?:\\000026)?(1Z[0-9A-Z]{16})$/);
        if (matches && matches.length > 1) {
            onBarcodeScan?.(scannerId, matches[1]);
            return true;
        }
    };

    const startInterval = () => {
        intervalRef.current = setInterval(() => {
            keyPresses.current = "";
        }, 50);
    };

    const handleKeyDown = (e: KeyboardEvent) => {
        if (e.key.length === 1) {
            keyPresses.current += e.key;
            clearInterval(intervalRef.current);
            startInterval();
            if (processKeypresses(undefined, keyPresses.current)) {
                keyPresses.current = "";
            }
        }
    };

    useEffect(() => {
        if (mountedRef.current) {
            return;
        }
        mountedRef.current = true;

        document.addEventListener("keydown", handleKeyDown);

        // timeout for keypresses
        startInterval();
    }, []);

    useEffect(() => {
        return () => {
            if (!mountedRef.current) {
                return;
            }
            mountedRef.current = false;
            clearInterval(intervalRef.current);
            document.removeEventListener("keydown", handleKeyDown);
        };
    }, []);

    useEffect(() => {
        if (barcodeData) {
            try {
                const msg: ScannerResponse = JSON.parse(barcodeData!);
                if (msg.messageType !== "barcode" || msg.barcode === undefined) {
                    return;
                }
                processKeypresses(msg.scannerId, msg.barcode)
            } catch (e) {
                if (e instanceof Error) {
                    toast.error(`Error parsing scanner response JSON: ${e.message}`);
                }
                toast.error("Unknown error parsing scanner response JSON");
            }
        }
    }, [barcodeData]);

    return <></>;
}
