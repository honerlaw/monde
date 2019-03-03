import * as React from "react";

interface IProps {
    label?: string;
    name: string;
    type: string;
    placeholder: string;
    help?: string;
    description?: string;

    inputProps?: { [key: string]: string };
}

export class InputGroup extends React.Component<IProps, {}> {

    private readonly inputId: string;
    private readonly helpId: string;

    public constructor(props: IProps) {
        super(props);

        this.inputId = this.getInputId();
        this.helpId = Math.random().toString().replace(".", "");
    }

    public render(): JSX.Element {
        return <div className="form-group">
            {this.props.label ? <label htmlFor={this.inputId}>{this.props.placeholder}</label> : null}
            <input className="form-control"
                   id={this.inputId}
                   name={this.props.name}
                   type={this.props.type}
                   aria-describedby={this.helpId}
                   placeholder={this.props.placeholder}
                   {...this.props.inputProps}/>
            {this.renderHelp()}
        </div>;
    }

    private renderHelp(): JSX.Element | null {
        if (!this.props.description) {
            return null;
        }
        return <small id={this.helpId} className="form-text text-muted">{this.props.description}</small>;
    }

    private getHelpId(): string {
        const id: string = Math.random().toString().replace(".", "");
        return `${id}-help`;
    }

    private getInputId(): string {
        const id: string = Math.random().toString().replace(".", "");
        return `${id}-input`;
    }

}
