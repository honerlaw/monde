import * as React from "react";
import {TextareaGroup} from "../../../bootstrap/textarea-group";

interface IProps {
    mediaId: string;
    parentCommentId?: string | null | undefined;

}

export class CommentForm extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <form className={"comment-form"} method={"POST"} action={`/media/comment/${this.props.mediaId}`}>
            {this.props.parentCommentId ? <input type={"hidden"} name={"parent_comment_id"} value={this.props.parentCommentId}/> : null}
            <TextareaGroup name={"comment"} placeholder={"add a comment"}/>
            <button type={"submit"}>post</button>
        </form>;
    }
}