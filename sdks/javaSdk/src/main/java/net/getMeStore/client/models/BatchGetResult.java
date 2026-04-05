
package net.getMeStore.client.models;

import java.util.List;
import java.util.Map;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class BatchGetResult {
    // key is the key that was found, value is the corresponding value
    private Map<String, String> found;
    
    // list of keys that were not found in the store
    private List<String> notFound;
    
    // key is the key that failed to get, value is the error message
    private Map<String, String> errors;
}