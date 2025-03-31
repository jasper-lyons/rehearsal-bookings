function handleBeforeUnload(event) {
if (startSlot || endSlot || selectedRoom) {
    event.preventDefault();
    event.returnValue = 'Are you sure you want to leave? Your changes may not be saved.';
}
}

window.beforeUnloadListenerAdded = false
