import React, { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";

interface ApiContextValue {
    apiUrl: string;
}

const ApiContext = createContext({} as ApiContextValue);

export const ApiProvider= ({ children }: {children: React.ReactNode}) => {
    const [apiUrl, setApiUrl] = useState<string>("");

    useEffect(() => {
        const fetchConfig = async () => {
            try {
                const response = await axios.get("/api/config.json");
                const config = await response.data;
                setApiUrl(config.apiUrl);
            } catch (error) {
                console.error("Error fetching config:", error);
            }
        }
        fetchConfig();
    }, []);


    return <ApiContext.Provider value={{ apiUrl }}> {children} </ApiContext.Provider>
    
}

// Custom hook to use the ApiContext
export const useApiContext = () => useContext(ApiContext);
