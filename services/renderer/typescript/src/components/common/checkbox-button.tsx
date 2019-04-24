import * as React from "react";
import "./checkbox-button.scss";

interface IProps {
    id: string;
    label: string;
}

export class CheckboxButton extends React.Component<IProps, {}> {

    private readonly id: string = Math.random().toString();

    public render(): JSX.Element | null {
        return <div className={"checkbox-button"}>
            <label htmlFor={`checkbox-button-toggle-${this.props.id}`}>{this.props.label}</label>
            <input type="checkbox" id={`checkbox-button-toggle-${this.props.id}`}/>
            <div className="content">
                {this.props.children}
            </div>
        </div>;
    }

}
