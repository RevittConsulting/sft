import { ColumnDef } from "@tanstack/react-table"
import { Toggle } from "@/types/toggles"



export const columns: ColumnDef<Toggle>[] = [
    {
        header: "Feature name",
        accessorKey: "feature_name"
    },
    {
        header: "Enabled",
        accessorKey: "enabled"
    },
    {
        header: "Toggle meta",
        accessorKey: "toggle_meta"
    },
    {
        header: "id",
        accessorKey: "id",
        cell: (id) => {
            return (
                <button onClick={() => alert(`${id.cell}`)}>Toggle</button>
            )
        }
    }
    // Probably abandon datatable and instead create a custom toggle display

]