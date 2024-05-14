import { ToggleDto } from "@/types/toggles";
import axios from "axios";

const toggles_url = "http://localhost:80/api/sft/v1/toggles"

export const fetchToggles = async () => {
    try {
        const response = await axios.get(`${toggles_url}`)
        console.log(response)
        return response
    } catch (error) {
        console.error("Error fetching toggles:", error)
    }
}

export const toggleFeature = async (toggleId: string) => {
    try {
        const response = await axios.patch(`${toggles_url}/${toggleId}`)
        return response
    } catch (error) {
        console.error("Error toggling feature:", error)
    }
}

export const createToggle = async (toggle: ToggleDto) => {
    try {
        const response = await axios.post(`${toggles_url}`, toggle)
        console.log(response)
        return response
    } catch (error) {
        console.error("Error creating toggle:", error)
    }
}

export const deleteToggle = async (toggleId: string) => {
    try {
        const response = await axios.delete(`${toggles_url}/${toggleId}`)
        console.log(response)
        return response
    } catch (error) {
        console.error("Error deleting toggle:", error)
    }
}