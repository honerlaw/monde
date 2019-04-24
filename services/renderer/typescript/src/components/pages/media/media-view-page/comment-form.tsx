import * as React from "react";
import {TextareaGroup} from "../../../bootstrap/textarea-group";

interface IProps {
    mediaId: string;
    parentCommentId?: string | null | undefined;
    hidden: boolean;
}

export class CommentForm extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <form className={"comment-form"}
                     method={"POST"}
                     action={`/media/comment/${this.props.mediaId}`}
                     style={{display: this.props.hidden ? 'none' : 'block'}}>
            {this.renderParentInput()}
            <TextareaGroup name={"comment"} placeholder={"add a comment"}/>
            <button className={"btn btn-primary"} type={"submit"}>post</button>
        </form>;
    }

    private renderParentInput(): JSX.Element | null {
        if (!this.props.parentCommentId) {
            return null;
        }
        return <input type={"hidden"} name={"parent_comment_id"} value={this.props.parentCommentId}/>;
    }

}
