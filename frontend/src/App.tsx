
import './App.css'
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useEffect, useState } from 'react';
import type { ScannerResponse } from './BarcodeScanner';
import { Shipping } from './generated/Shipping';
import type { ShippingResponse } from './generated/data-contracts';

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
    const [barcodeScannerImage, setBarcodeScannerImage] = useState("")

    // Process image data from barcode scanner through websocket
    useEffect(() => {
        if (readyState === ReadyState.OPEN && lastMessage !== null) {
            try {
                const msg: ScannerResponse = JSON.parse(lastMessage.data);
                if (msg.messageType === "image") {
                    setBarcodeScannerImage(msg.image!);
                    shippingApi.labelCheckCreate({ image: msg.image })
                        .then(({ data }: { data: ShippingResponse }) => {
                            console.log("Label check response:", data);
                            // TODO: Show red/green based on label validity
                        }).catch((error) => {
                            console.error("Error during label check:", error);
                        });
                }
            } catch (e) {
                if (e instanceof Error) {
                    console.error(`Error parsing scanner response JSON: ${e.message}`);
                } else {
                    console.error("Unknown error parsing scanner response JSON");
                }
            }
        }
    }, [lastMessage]);

    return (
        <>
            <div className="card">
                <button
                    onClick={() => sendJsonMessage(
                        JSON.stringify({
                            scannerId: 1,
                            commandType: "image_mode",
                        }))}
                    disabled={readyState !== ReadyState.OPEN}
                >Take Image</button>
                {barcodeScannerImage && (
                    <img src={`data:image/jpeg;charset=utf-8;base64,${barcodeScannerImage}`} />
                )}
            </div >
        </>
    )
}

export default App
