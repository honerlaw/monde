import * as React from "react";

interface IProps {
    label?: string;
    name: string;
    placeholder: string;
    value?: string;
    help?: string;
    description?: string;

    textareaProps?: { [key: string]: string };
}

export class TextareaGroup extends React.Component<IProps, {}> {

    private readonly id: string;
    private readonly helpId: string;

    public constructor(props: IProps) {
        super(props);

        this.id = this.getId();
        this.helpId = Math.random().toString().replace(".", "");
    }

    public render(): JSX.Element {
        return <div className="form-group">
            {this.props.label ? <label htmlFor={this.id}>{this.props.placeholder}</label> : null}
            <textarea className="form-control"
                   id={this.id}
                   name={this.props.name}
                   aria-describedby={this.helpId}
                   placeholder={this.props.placeholder}
                      {...this.props.textareaProps}>
                {this.props.value}
            </textarea>
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

    private getId(): string {
        const id: string = Math.random().toString().replace(".", "");
        return `${id}-textarea`;
    }

}
