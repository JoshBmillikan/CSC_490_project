/** @jsxImportSource @emotion/react */
import {css} from "@emotion/react";
import {useEffect, useState} from "react";
import {useAppDispatch, useAppSelector} from "./data/store";
import {ShaderState} from "./data/reducers";

export function CodeEditor() {
    const selector = useAppSelector((state) => state.shader)
    const [getText, setText] = useState(selector[selector.selectedShaderName as keyof ShaderState] ?? "")
    const dispatch = useAppDispatch()
    useEffect(()=> {
        const t = setTimeout(()=>{
            dispatch({shaderName: selector.selectedShaderName, newShader: getText, type: "updateShader"})
        },1000)
        return () => clearTimeout(t)
    },[dispatch, getText, selector.selectedShaderName])

    return (
        <div css={css`
          padding-left: 10%;
          padding-top: 10%;
          user-select: none;
        `}>
        <textarea
            value={getText}
            onChange={(event => setText(event.target.value))}
            css={theme => ({
                height: '80vh',
                width: '50vh',
                backgroundColor: theme.foregroundColor,
                color: theme.textColor,
                resize: 'none',
                borderStyle: 'solid',
                borderColor: theme.borderColor,
                borderRadius: "3%",
                padding: '10px'
            })}
        />
        </div>
    )
}