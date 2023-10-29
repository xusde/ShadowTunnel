import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {GetConfigValue, SetConfig} from "../../wailsjs/go/main/App";

export interface SettingsState {
    mode: string;
    localPort: string;
    proxyAddress: string;
    encryptionMethod: string;
    encryptionKey: string;
    transportProtocol: string;
}

const initialState: SettingsState = {
    mode: "",
    localPort: "",
    proxyAddress: "",
    encryptionMethod: "",
    encryptionKey: "",
    transportProtocol: "",
}

export const settingsSlice = createSlice({
    name: 'settings',
    initialState,
    reducers: {
        setMode: (state, action:PayloadAction<string>) => {
            state.mode = action.payload;
            SetConfig(marshalSettings(state));
        },
        setLocalPort: (state, action:PayloadAction<string>) => {
            state.localPort = action.payload;
            SetConfig(marshalSettings(state));
        },
        setProxyAddress: (state, action:PayloadAction<string>) => {
            state.proxyAddress = action.payload;
            SetConfig(marshalSettings(state));
        },
        setEncryptionMethod: (state, action:PayloadAction<string>) => {
            state.encryptionMethod = action.payload;
            SetConfig(marshalSettings(state));
        },
        setEncryptionKey: (state, action:PayloadAction<string>) => {
            state.encryptionKey = action.payload;
            SetConfig(marshalSettings(state));
        },
        setTransportProtocol: (state, action:PayloadAction<string>) => {
            state.transportProtocol = action.payload;
            SetConfig(marshalSettings(state));
        }
    }
})

function marshalSettings(settings: SettingsState) {
    return JSON.stringify(settings);
}

export const {setMode, setLocalPort, setProxyAddress, setEncryptionMethod, setEncryptionKey, setTransportProtocol} = settingsSlice.actions;

export default settingsSlice.reducer;