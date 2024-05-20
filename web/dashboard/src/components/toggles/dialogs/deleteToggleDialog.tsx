import { Button } from "../../ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
  } from "../../ui/dialog";

interface DeleteToggleDialogProps {
    isOpen: boolean;
    setIsOpen: (isOpen: boolean) => void;
    deleteToggle: () => void;
}

export function DeleteToggleDialog({
    isOpen,
    setIsOpen,
    deleteToggle,
}: DeleteToggleDialogProps) {
    return (
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete task</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this task?
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <div className="w-full flex justify-center gap-4">
              <Button type="button" variant="outline" onClick={() => setIsOpen(false)}>
                Cancel
              </Button>
              <Button type="button" onClick={() => {
                deleteToggle()
                setIsOpen(false)
                }}>
                Delete
              </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
    )
}