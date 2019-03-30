import * as React from "react";

interface IProps {
    label?: string;
    name: string;
    type: string;
    placeholder: string;
    value?: string;
    help?: string;
    description?: string;

    inputProps?: { [key: string]: string };
}

export class InputGroup extends React.Component<IProps, {}> {

    private readonly inputId: string;
    private readonly helpId: string;

    public constructor(props: IProps) {
        super(props);

        this.inputId = `${this.props.name}-input`;
        this.helpId = `${this.props.name}-help`;
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
                   value={this.props.value ? this.props.value : undefined}
                   defaultValue={this.props.value ? undefined : this.props.value}
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

}
