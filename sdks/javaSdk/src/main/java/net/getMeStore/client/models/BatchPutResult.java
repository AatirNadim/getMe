
package net.getMeStore.client.models;

import java.util.List;
import java.util.Map;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;


@Data
@NoArgsConstructor
@AllArgsConstructor
public class BatchPutResult {
    // list of keys that were successfully put in the store
    private Integer successful;
    
    // key is the key that failed to put, value is the error message
    private Map<String, String> errors;
}