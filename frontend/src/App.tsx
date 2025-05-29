
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useEffect, useState } from 'react';
import type { ScannerResponse } from './BarcodeScanner';
import { Shipping } from './generated/Shipping';
import type { ValidationResponse, ValidationResult } from './generated/data-contracts';
import { ClipLoader } from "react-spinners";
import { toast, ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import BarcodeListener from "./BarcodeListener";

function App() {
    const { VITE_API_BASE_URL } = import.meta.env;
    const shippingApi = new Shipping({
        baseURL: `${VITE_API_BASE_URL}`,
    });
    const {
        lastMessage,
        sendJsonMessage,
        readyState,
    } = useWebSocket("ws://172.20.10.2:8089", {
        shouldReconnect: () => true, // Always try to reconnect
        reconnectAttempts: 100000, // Always retry
        reconnectInterval: 3 * 1000, // Reconnect attempt interval in milliseconds
    });
    const [scanState, setScanState] = useState<"idle" | "validating" | "invalid">("idle"); // Possible states: idle, scanning, processing
    const [barcode, setBarcode] = useState("");
    const [barcodeScannerImage, setBarcodeScannerImage] = useState("");
    const [result, setResult] = useState<ValidationResult>();

    const successBeep = new Audio("/beep-success.mp3");
    const errorBeep = new Audio("/beep-error.mp3");

    // Process image data from barcode scanner through websocket
    useEffect(() => {
        if (readyState === ReadyState.OPEN && lastMessage !== null) {
            try {
                const msg: ScannerResponse = JSON.parse(lastMessage.data);
                if (msg.messageType === "image") {
                    setScanState("validating");
                    setBarcodeScannerImage(msg.image!);
                    shippingApi.labelValidateCreate({
                        trackingNumber: barcode,
                        image: msg.image,
                    })
                        .then(({ data }: { data: ValidationResponse }) => {
                            setResult(data.result);
                            if (data.result?.valid) {
                                successBeep.play();
                                toast.success("Valid. Scan next label.");
                                setScanState("idle");
                            } else {
                                errorBeep.play();
                                toast.error("Invalid label. Reprint label.");
                                setScanState("invalid");
                            }
                        }).catch((error) => {
                            toast.error("Error validating:", error);
                        });
                }
            } catch (e) {
                if (e instanceof Error) {
                    toast.error(`Error parsing scanner response JSON: ${e.message}`);
                } else {
                    toast.error("Unknown error parsing scanner response JSON");
                }
            }
        }
    }, [lastMessage]);

    const onBarcodeScan = (scannerId: number | undefined, barcode: string) => {
        setBarcode(barcode);
        sendJsonMessage(
            JSON.stringify({
                scannerId: scannerId,
                commandType: "image_mode",
            })
        );
    }

    const getBody = () => {
        switch (scanState) {
            case "idle":
                return (
                    <div
                        className="flex flex-col items-center justify-center border-2 border-black rounded-[40px] p-16 gap-4">
                        <div className="text-8xl text-[#301506]">Scan a label</div>
                        <div className="text-2xl text-black">From a foot away, so you get all of it</div>
                        <img src="/scan-label.svg" alt="Scan Label" className="w-64 h-64 mt-8" />
                    </div>
                );
            case "validating":
                return (
                    <div className="flex flex-col items-center justify-center gap-8">
                        <div className="text-2xl">Validating</div>
                        <ClipLoader
                            color="#000000"
                            loading={true}
                            size={150}
                            aria-label="Loading Spinner"
                            data-testid="loader"
                        />
                    </div>
                );
            case "invalid":
                return (
                    <div className="flex flex-col gap-16 items-center">
                        <div className="flex flex-row gap-32">
                            <div className="flex flex-col items-center justify-center border-2 border-black rounded-[40px] p-8">
                                <div className="w-full text-start text-xl">Scanned</div>
                                <img src={`data:image/jpeg;charset=utf-8;base64,${barcodeScannerImage}`} alt="Scanned Label" className="w-96 h-96" />
                            </div>
                            <div className="flex flex-col items-center justify-center border-2 border-black rounded-[40px] p-8">
                                <div className="w-full text-start text-xl">Expected</div>
                                {result?.expectedAddress?.addressLine1}<br />
                                {result?.expectedAddress?.addressLine2 && `${result?.expectedAddress?.addressLine2}<br />`}
                                {result?.expectedAddress?.city}, {result?.expectedAddress?.stateProvince} {result?.expectedAddress?.postalCode}
                            </div>
                        </div>
                        <div className="w-64 p-4 border-4 border-[#301506] bg-[#301506] rounded-lg text-4xl text-[#FAB80A] cursor-pointer text-center" onClick={() => {
                            toast.success("Reprinted. Reprint label.");
                            setScanState("idle");
                        }}>Reprint label</div>
                    </div>
                );
        }
    }

    return (
        <>
            <div className="flex flex-col h-screen w-full">
                <div className="flex flex-row bg-[#301506] items-center">
                    <img src="/ups.svg" alt="UPS Logo" className="w-16 h-16 p-4" />
                    <div className="text-2xl text-[#FAB80A]">Validate</div>
                </div>
                <div className="flex flex-col h-full w-full items-center justify-center">
                    {getBody()}
                </div>
            </div>
            <BarcodeListener barcodeData={lastMessage?.data} onBarcodeScan={onBarcodeScan} />
            <ToastContainer
                position="top-right"
                autoClose={3000}
                hideProgressBar={false}
                newestOnTop={false}
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
            />
        </>
    )
}

export default App
