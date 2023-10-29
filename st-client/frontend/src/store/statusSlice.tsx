import {createSlice, PayloadAction} from "@reduxjs/toolkit";

export interface StatusState {
    isConnected: boolean;
}

const initialState: StatusState = {
    isConnected: false,
}

export const statusSlice = createSlice({
    name: 'status',
    initialState,
    reducers: {
        setConnected: (state, action:PayloadAction<boolean>) => {
            state.isConnected = action.payload;
        },
        toggleConnection: (state) => {
            state.isConnected = !state.isConnected;
        }
    }
})

export const {setConnected, toggleConnection} = statusSlice.actions;

export default statusSlice.reducer;