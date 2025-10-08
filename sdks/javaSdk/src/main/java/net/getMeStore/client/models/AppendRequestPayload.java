package net.getMeStore.client.models;


import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;


@Getter
@Setter
public class AppendRequestPayload {

    @JsonProperty("Key")
    String Key;
    @JsonProperty("Value")
    String Value;


    public AppendRequestPayload(String Key, String Value) {
        this.Key = Key;
        this.Value = Value;
    }

}
