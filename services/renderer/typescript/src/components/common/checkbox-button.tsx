import * as React from "react";
import "./checkbox-button.scss";

interface IProps {
    id: string;
    label: string;
}

export class CheckboxButton extends React.Component<IProps, {}> {

    private readonly id: string = Math.random().toString();

    public render(): JSX.Element | null {
        const children: React.ReactNode[] = React.Children.toArray(this.props.children);
        return <div className={"checkbox-button"}>
            <label className={"checkbox-button-label"} htmlFor={`checkbox-button-toggle-${this.props.id}`}>{this.props.label}</label>
            <input type="checkbox" id={`checkbox-button-toggle-${this.props.id}`}/>
            <div className={"content-hidden"}>{children[0]}</div>
            <div className={"content-visible"}>{children[1]}</div>
        </div>;
    }

}
